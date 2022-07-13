package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/enclaive/relay/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
	"time"
)

func (s *Server) EnclaveCreator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("x-username")
			lookup, err := s.repoManager.Lookup().GetByUsername(c.Request().Context(), username)

			if s.repoManager.IsEmptyResultSetError(err) {
				log.Info().Caller().Msg("spawning new enclave")
				ip, err := s.DeployEnclave(c.Request().Context())
				if err != nil {
					log.Error().Caller().Err(err).Msg("failed to spawn enclave")
					return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				lookup = models.Lookup{
					Username:       username,
					EnclaveAddress: ip,
				}

				err = s.repoManager.Lookup().Set(c.Request().Context(), lookup)
				if err != nil {
					log.Error().Caller().Err(err).Msg("failed to save enclave ip")
					return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}

				log.Info().Caller().Msgf("registered new enclave: %v", lookup)

				c.Set("address", lookup.EnclaveAddress)

				time.Sleep(10 * time.Second)

				return next(c)

			} else if err != nil {
				log.Error().Caller().Err(err).Msg("failed to get lookup from database")
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			} else {
				log.Info().Caller().Err(err).Msg("user tried to register an already used account")
				return c.String(http.StatusBadRequest, "username is already taken")
			}
		}
	}
}

func (s *Server) DeployEnclave(ctx context.Context) (string, error) {
	appDeploymentsClient := s.clientset.AppsV1().Deployments("default")
	secret, _ := s.clientset.CoreV1().Secrets("default").Get(ctx, "regcred", metav1.GetOptions{})
	backendip, _ := s.clientset.CoreV1().Services("default").Get(ctx, "backend-service", metav1.GetOptions{})
	randsubstr, _ := randomHex(16)
	randsubstr = "enclave" + randsubstr
	quantity, quantityErr := resource.ParseQuantity("512Ki")
	if quantityErr != nil {
		log.Error().Caller().Err(quantityErr).Msg("failed to parse quantity")
	}
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
							/*							VolumeMounts: []apiv1.VolumeMount{
														{
															Name:      "user-",
															MountPath: "/server",
														},
														{
															Name:      "sgx-volume",
															MountPath: "/dev/sgx",
														},
													},*/
							Env: []apiv1.EnvVar{
								{
									Name:  "BACKEND_IP",
									Value: fmt.Sprintf("%s:%d", backendip.Spec.ClusterIP, backendip.Spec.Ports[0].Port),
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 2533,
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
					/*					Volumes: []apiv1.Volume{
										{
											Name: "server-volume",
											VolumeSource: apiv1.VolumeSource{
												EmptyDir: &apiv1.EmptyDirVolumeSource{},
											},
										},
										{
											Name: "sgx-volume",
											VolumeSource: apiv1.VolumeSource{
												EmptyDir: &apiv1.EmptyDirVolumeSource{}},
										},
									},*/
				},
			},
		},
	}

	appResult, err := appDeploymentsClient.Create(ctx, appDeployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	serviceDeploymentsClient := s.clientset.CoreV1().Services("default")

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
					TargetPort: intstr.FromInt(2533),
				},
			},
		},
	}

	serviceResult, err := serviceDeploymentsClient.Create(ctx, serviceDeployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	service, err := s.clientset.CoreV1().Services("default").Get(ctx, serviceResult.GetObjectMeta().GetName(), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, service.Spec.Ports[0].Port), nil
}

func int32Ptr(i int32) *int32 { return &i }

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
