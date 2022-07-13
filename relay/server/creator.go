package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/enclaive/relay/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"strings"
	"time"
)

const (
	NAMESPACE = "default"
	TIMEOUT   = 180 * time.Second
)

func (s *Server) EnclaveCreator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("x-username")
			lookup, err := s.repoManager.Lookup().GetByUsername(c.Request().Context(), username)

			if s.repoManager.IsEmptyResultSetError(err) {
				log.Info().Caller().Msg("spawning new enclave")
				ip, podName, err := s.DeployEnclave(c.Request().Context())
				if err != nil {
					log.Error().Caller().Err(err).Msg("failed to spawn enclave")
					return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				lookup = models.Lookup{
					Username:       username,
					EnclaveAddress: ip,
				}

				err = s.repoManager.Lookup().Set(c.Request().Context(), lookup)
				if err != nil {
					log.Error().Caller().Err(err).Msg("failed to save enclave ip")
					return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				log.Info().Caller().Msgf("registered new enclave: %v", lookup)

				c.Set("address", lookup.EnclaveAddress)

				//time.Sleep(20 * time.Second)
				if err := waitForPodRunning(c.Request().Context(), s.clientset, podName); err != nil {
					log.Error().Caller().Err(err).Msg("enclave failed to start")
					return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				return next(c)

			} else if err != nil {
				log.Error().Caller().Err(err).Msg("failed to get lookup from database")
				return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			} else {
				log.Info().Caller().Err(err).Msg("user tried to register an already used account")
				return echo.NewHTTPError(http.StatusBadRequest, "username is already taken")
			}
		}
	}
}

func (s *Server) DeployEnclave(ctx context.Context) (string, string, error) {
	appDeploymentsClient := s.clientset.AppsV1().Deployments(NAMESPACE)
	secret, _ := s.clientset.CoreV1().Secrets(NAMESPACE).Get(ctx, "regcred", metav1.GetOptions{})
	backendip, _ := s.clientset.CoreV1().Services(NAMESPACE).Get(ctx, "backend-service", metav1.GetOptions{})
	randsubstr, _ := randomHex(16)
	randsubstr = "enclave" + randsubstr
	quantity, _ := resource.ParseQuantity("512Ki")

	appDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": randsubstr,
			},
			GenerateName: "enclave-",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":   randsubstr,
					"tier":  "backend",
					"track": "stable",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":   randsubstr,
						"tier":  "backend",
						"track": "stable",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "enclave",
							Image: s.cfg.Image,
							Env: []apiv1.EnvVar{
								{
									Name:  "BACKEND_IP",
									Value: fmt.Sprintf("%s:%d", backendip.Spec.ClusterIP, backendip.Spec.Ports[0].Port),
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									"sgx.intel.com/epc": quantity,
								},
							},
						},
					},
					ImagePullSecrets: []apiv1.LocalObjectReference{{secret.Name}},
					NodeSelector: map[string]string{
						"disktype": "ssd",
					},
				},
			},
		},
	}

	appResult, err := appDeploymentsClient.Create(ctx, appDeployment, metav1.CreateOptions{})
	if err != nil {
		return "", "", err
	}
	serviceDeploymentsClient := s.clientset.CoreV1().Services(NAMESPACE)

	serviceDeployment := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-service", appResult.Name),
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": randsubstr,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       2534,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	serviceResult, err := serviceDeploymentsClient.Create(ctx, serviceDeployment, metav1.CreateOptions{})
	if err != nil {
		return "", "", err
	}

	service, err := s.clientset.CoreV1().Services(NAMESPACE).Get(ctx, serviceResult.GetObjectMeta().GetName(), metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, service.Spec.Ports[0].Port), appResult.GetObjectMeta().GetName(), nil
}

func int32Ptr(i int32) *int32 { return &i }

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// return a condition function that indicates whether the given pod is
// currently running
func isPodRunning(ctx context.Context, c kubernetes.Interface, podName string) wait.ConditionFunc {
	return func() (bool, error) {
		pods, err := c.CoreV1().Pods(NAMESPACE).List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.GetName(), podName) {
				fmt.Println("found match")
				switch pod.Status.Phase {
				case apiv1.PodRunning:
					return true, nil
				default:
					return false, nil
				}
			}
		}

		return false, errors.New("pod does not exist")
	}
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func waitForPodRunning(ctx context.Context, c kubernetes.Interface, podName string) error {
	return wait.PollImmediate(time.Second, TIMEOUT, isPodRunning(ctx, c, podName))
}
