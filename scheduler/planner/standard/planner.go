package standard

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"log"
)

type StandardPlanner struct {
	plan *Plan
}

func (s StandardPlanner) GeneratePlan(configuration *configuration.Configuration,
	flag configuration.CommandFlag) (planner.IPlan, error) {

	s.plan = NewPlan(flag)

	//Clusters -----------------------------------------------------------------------------------

	for _, clusterKinds := range configuration.ClusterKinds {

		for _, providers := range clusterKinds.Strategy {

			pcr := ProviderClusterReference{ProviderName: providers.Provider.Name}

			s.plan.providerClusterReference[providers.Provider.Name] = &pcr

			for _, clusters := range providers.Provider.Clusters {
				s.plan.clusterDeployments[pcr.ProviderName] = append(s.plan.clusterDeployments[pcr.ProviderName], clusters.Cluster)

				pcr.Dependencies = append(pcr.Dependencies, clusters.Cluster)
			}

		}
	}
	log.Printf(fmt.Sprintf("%#v", s.plan.clusterDeployments))
	// -------------------------------------------------------------------------------------------

	// Application -------------------------------------------------------------------------------
	for _, applicationKinds := range configuration.ApplicationKinds {
		for _, applicationLists := range applicationKinds.Strategy {
			//Name of cluster (short name)

			for _, app := range applicationLists.Cluster.Applications {

				s.plan.clusterApplications[applicationLists.Cluster.ShortName] =
					append(s.plan.clusterApplications[applicationLists.Cluster.ShortName], app.Application)
			}
		}
	}
	// -------------------------------------------------------------------------------------------
	log.Printf(fmt.Sprintf("%#v", s.plan.clusterApplications))
	// -------------------------------------------------------------------------------------------
	return s.plan, nil
}
