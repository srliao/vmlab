package defaults

import (
	vs1alpha1 "github.com/backube/volsync/api/v1alpha1"
	"github.com/srliao/vmlab/pkg/klusterhelper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// this file containers helper method to create a standard
// pvc synced with volsync
// external secret keys are hard coded

func NewVolsyncPVCResources(app, namespace, size string) []klusterhelper.KubeResource {
	resources := []klusterhelper.KubeResource{
		newVolsyncPVC(app, namespace, size),
		newBootstrap(app, size),
		newVolsyncSource(app, size),
		newVolsyncES(app),
	}
	return resources
}

func newVolsyncPVC(app, namespace, size string) klusterhelper.KubeResource {
	pvc := NewPVC(app, namespace, size).
		WithDataSourceRef(
			app+"-bootstrap",
			"ReplicationDestination",
			"volsync.backube",
		)
	return pvc
}

func newBootstrap(app, size string) klusterhelper.KubeResource {
	storageClass := defaultPVCStorageClass
	resSize := resource.MustParse(size)
	uid := defaultUID
	bootstrap := &klusterhelper.ReplicationDestinationWrapper{
		ReplicationDestination: &vs1alpha1.ReplicationDestination{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "volsync.backube/v1alpha1",
				Kind:       "ReplicationDestination",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: app + "-bootstrap",
			},
			Spec: vs1alpha1.ReplicationDestinationSpec{
				Trigger: &vs1alpha1.ReplicationDestinationTriggerSpec{
					Manual: "restore-once-bootstrap",
				},
				Restic: &vs1alpha1.ReplicationDestinationResticSpec{
					Repository:            app + "-volsync-minio",
					CacheStorageClassName: &storageClass,
					CacheCapacity:         &resSize,
					CacheAccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					ReplicationDestinationVolumeOptions: vs1alpha1.ReplicationDestinationVolumeOptions{
						CopyMethod:              vs1alpha1.CopyMethodSnapshot,
						VolumeSnapshotClassName: &storageClass,
						StorageClassName:        &storageClass,
						Capacity:                &resSize,
						AccessModes:             []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					},
					MoverConfig: vs1alpha1.MoverConfig{
						MoverSecurityContext: &corev1.PodSecurityContext{
							RunAsUser:  &uid,
							RunAsGroup: &uid,
							FSGroup:    &uid,
						},
					},
				},
			},
		},
	}
	return bootstrap
}

func newVolsyncES(app string) klusterhelper.KubeResource {
	return NewESWithDataAndKey(
		app+"-volsync-minio",
		"volsync",
		map[string]string{
			"RESTIC_REPOSITORY":     "s3:http://minio.lan/persistent-snapshots/" + app,
			"RESTIC_PASSWORD":       "{{ .RESTIC_PASSWORD }}",
			"AWS_ACCESS_KEY_ID":     "{{ .AWS_ACCESS_KEY_ID }}",
			"AWS_SECRET_ACCESS_KEY": "{{ .AWS_SECRET_ACCESS_KEY }}",
		},
		"volsync",
	)
}

func newVolsyncSource(app, size string) klusterhelper.KubeResource {
	storageClass := defaultPVCStorageClass
	resSize := resource.MustParse(size)
	uid := defaultUID
	schedule := "0 * * * *" // every hour
	source := &klusterhelper.ReplicationSourceWrapper{
		ReplicationSource: &vs1alpha1.ReplicationSource{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "volsync.backube/v1alpha1",
				Kind:       "ReplicationSource",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: app + "-minio",
			},
			Spec: vs1alpha1.ReplicationSourceSpec{
				SourcePVC: app,
				Trigger: &vs1alpha1.ReplicationSourceTriggerSpec{
					Schedule: &schedule,
				},
				Restic: &vs1alpha1.ReplicationSourceResticSpec{
					Repository:            app + "-volsync-minio",
					CacheStorageClassName: &storageClass,
					CacheCapacity:         &resSize,
					CacheAccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					ReplicationSourceVolumeOptions: vs1alpha1.ReplicationSourceVolumeOptions{
						CopyMethod:              vs1alpha1.CopyMethodSnapshot,
						StorageClassName:        &storageClass,
						VolumeSnapshotClassName: &storageClass,
						AccessModes:             []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					},
					MoverConfig: vs1alpha1.MoverConfig{
						MoverSecurityContext: &corev1.PodSecurityContext{
							RunAsUser:  &uid,
							RunAsGroup: &uid,
							FSGroup:    &uid,
						},
					},
					PruneIntervalDays: klusterhelper.Int32Ptr(7),
					Retain: &vs1alpha1.ResticRetainPolicy{
						Hourly: klusterhelper.Int32Ptr(24),
						Daily:  klusterhelper.Int32Ptr(7),
						Weekly: klusterhelper.Int32Ptr(5),
					},
				},
			},
		},
	}
	return source
}
