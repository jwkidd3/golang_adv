package gameserver

import (
	"context"

	gameserverv1 "github.com/motiso/gameserveroperator/pkg/apis/gameserver/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const backendPort = 8080
const backendServicePort = 30685
const backendImage = "motisoffer/gogameserver:1.0"

func backendDeploymentName(v *gameserverv1.Gameserver) string {
	return v.Name + "-backend"
}

func backendServiceName(v *gameserverv1.Gameserver) string {
	return v.Name + "-backend-service"
}

func configMapName(v *gameserverv1.Gameserver) string {
	return v.Name + "-config"
}
func (r *ReconcileGameserver) backendDeployment(v *gameserverv1.Gameserver) *appsv1.Deployment {
	labels := labels(v, "backend")
	// we dont want more then one pn purpose
	size := int32(1)

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
			Key:                  "password",
		},
	}

	id := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: configMapName(v)},
			Key:                  "GAME_ID",
		},
	}
	name := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: configMapName(v)},
			Key:                  "GAME_NAME",
		},
	}
	description := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: configMapName(v)},
			Key:                  "GAME_DESCRIPTION",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           backendImage,
							ImagePullPolicy: corev1.PullAlways,
							Name:            "gameserver",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: backendPort,
									Name:          "gameserver",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_DATABASE",
									Value: "users_db",
								},
								{
									Name:  "MYSQL_SERVICE_HOST",
									Value: mysqlServiceName(),
								},
								{
									Name:      "MYSQL_USERNAME",
									ValueFrom: userSecret,
								},
								{
									Name:      "MYSQL_PASSWORD",
									ValueFrom: passwordSecret,
								},
								{
									Name:      "GAME_ID",
									ValueFrom: id,
								},
								{
									Name:      "GAME_NAME",
									ValueFrom: name,
								},
								{
									Name:      "GAME_DESCRIPTION",
									ValueFrom: description,
								},
							},
						}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcileGameserver) backendService(v *gameserverv1.Gameserver) *corev1.Service {
	labels := labels(v, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendServiceName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
				NodePort:   v.Spec.ServerPort,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	log.Info("--------------------------node port is " + v.Spec.Name)
	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcileGameserver) updateBackendStatus(v *gameserverv1.Gameserver) error {
	//v.Status.BackendImage = backendImage
	err := r.client.Status().Update(context.TODO(), v)
	return err
}

func (r *ReconcileGameserver) handleBackendChanges(v *gameserverv1.Gameserver) (*reconcile.Result, error) {

	found := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: configMapName(v), Namespace: v.Namespace}, found)

	id := v.Spec.GameID
	name := v.Spec.Name
	description := v.Spec.Description

	if id != found.Data["GAME_ID"] ||
		name != found.Data["GAME_NAME"] ||
		description != found.Data["GAME_DESCRIPTION"] {

		found.Data["GAME_ID"] = id
		found.Data["GAME_NAME"] = name
		found.Data["GAME_DESCRIPTION"] = description

		err = r.client.Update(context.TODO(), found)
		if err != nil {
			log.Error(err, "Failed to update ConfigMap.", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
			return &reconcile.Result{}, err
		}

		//delete the dependent deployment and requeue
		found := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      backendDeploymentName(v),
			Namespace: v.Namespace,
		}, found)

		err = r.client.Delete(context.TODO(), found)
		if err != nil {
			return &reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return &reconcile.Result{Requeue: true}, nil
	}

	return nil, nil
}

func (r *ReconcileGameserver) newConfigMap(v *gameserverv1.Gameserver) *corev1.ConfigMap {
	labels := map[string]string{
		"app": v.Name,
	}
	s := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(v),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"GAME_ID":          v.Spec.GameID,
			"GAME_NAME":        v.Spec.Name,
			"GAME_DESCRIPTION": v.Spec.Description,
		},
	}
	controllerutil.SetControllerReference(v, s, r.scheme)

	return s
}
