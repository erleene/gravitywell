package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1ServiceResouce(k kubernetes.Interface, objdep *v1.Service, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found service resource")
	ssclient := k.CoreV1().Services(namespace)

	if opts.DryRun {
		_, err := ssclient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Service resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Service resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = ssclient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := ssclient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			log.Debug(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		_, err := ssclient.Create(objdep)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Service resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Service deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := ssclient.Create(objdep)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Service resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Service deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := ssclient.Update(objdep)
		if err != nil {
			log.Error("Could not update Service")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Service updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := ssclient.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
