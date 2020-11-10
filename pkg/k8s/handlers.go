package k8s

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateAudioJob(c *gin.Context) {
	client := GetClient()
	var parall = int32(1)
	var comple = int32(1)
	jobClient := client.BatchV1().Jobs(apiv1.NamespaceDefault)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "meeting1",
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parall,
			Completions: &comple,
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "audiojob",
							Image: "harbor.beautytiger.com/docker.io/nginx:1.15.12",
						},
					},
				},
			},
		},
	}
	fmt.Println("create job")
	result, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		c.String(400, err.Error())
	} else {
		c.JSON(200, result)
	}

	return
}
