package server

import (
	"context"
	"github.com/enclaive/relay/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func (s *Server) UserAddressMapper() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("x-username")
			lookup, err := s.repoManager.Lookup().GetByUsername(c.Request().Context(), username)

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
	deploymentsClient := s.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	//TODO replace with real deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	pod, err := s.clientset.CoreV1().Pods("enclaves").Get(ctx, result.GetObjectMeta().GetName(), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pod.Status.PodIP, nil
}

func int32Ptr(i int32) *int32 { return &i }
