package server

import (
	"context"
	"fmt"
	"github.com/enclaive/relay/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
)

func (s *Server) UserAddressMapper() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("x-username")
			lookup, err := s.repoManager.Lookup().GetByUsername(c.Request().Context(), username)
			log.Log().Caller().Msg("received msg")

			if s.repoManager.IsEmptyResultSetError(err) {
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
			} else if err != nil {
				log.Error().Caller().Err(err).Msg("failed to get lookup from database")
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			c.Set("address", lookup.EnclaveAddress)

			return next(c)
		}
	}
}

func (s *Server) DeployEnclave(ctx context.Context) (string, error) {
	appDeploymentsClient := s.clientset.AppsV1().Deployments("enclave-ns")

	appDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "enclave",
			},
			GenerateName: "enclave-",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":   "enclave",
					"tier":  "backend",
					"track": "stable",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":   "enclave",
						"tier":  "backend",
						"track": "stable",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "enclave",
							Image: s.cfg.Image,
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 2533,
								},
							},
						},
					},
					ImagePullSecrets: []apiv1.LocalObjectReference{{"regcred"}},
					NodeSelector: map[string]string{
						"disktype": "ssd",
					},
				},
			},
		},
	}

	appResult, err := appDeploymentsClient.Create(ctx, appDeployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	serviceDeploymentsClient := s.clientset.CoreV1().Services("enclave-ns")

	serviceDeployment := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-service", appResult.Name),
			Namespace: "enclave-ns",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": appResult.Name,
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

	service, err := s.clientset.CoreV1().Services("enclave-ns").Get(ctx, serviceResult.GetObjectMeta().GetName(), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, service.Spec.Ports[0].Port), nil
}

func int32Ptr(i int32) *int32 { return &i }
