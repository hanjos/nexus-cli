package gavs // import "sbrubbles.org/go/nexus-cli/gavs"

import (
	"sbrubbles.org/go/nexus"

	"fmt"
)

type Gav struct {
	GroupID    string
	ArtifactID string
	Version    string
	Artifacts  []*nexus.Artifact
}

func newGav(g, a, v string) *Gav {
	return &Gav{GroupID: g, ArtifactID: a, Version: v, Artifacts: []*nexus.Artifact{}}
}

func (gav *Gav) hash() string {
	return fmt.Sprintf("%v:%v:%v", gav.GroupID, gav.ArtifactID, gav.Version)
}

func (gav *Gav) add(artifact *nexus.Artifact) {
	gav.Artifacts = append(gav.Artifacts, artifact)
}

func (gav *Gav) String() string {
	return gav.hash()
}

type gavSet struct {
	Data []*Gav
	Map  map[string]*Gav
}

func newGavSet() *gavSet {
	return &gavSet{
		Data: []*Gav{},
		Map:  map[string]*Gav{},
	}
}

func (set *gavSet) getGavOf(artifact *nexus.Artifact) *Gav {
	if gav, ok := set.Map[hashFromArtifact(artifact)]; ok {
		return gav
	}

	gav := newGav(artifact.GroupID, artifact.ArtifactID, artifact.Version)
	set.add(gav)

	return gav
}

func (set *gavSet) add(gav *Gav) {
	hash := gav.hash()

	if _, ok := set.Map[hash]; ok {
		return
	}

	set.Data = append(set.Data, gav)
	set.Map[hash] = gav
}

func hashFromArtifact(artifact *nexus.Artifact) string {
	return fmt.Sprintf("%v:%v:%v", artifact.GroupID, artifact.ArtifactID, artifact.Version)
}

func GavsOf(artifacts []*nexus.Artifact) []*Gav {
	set := newGavSet()
	for _, artifact := range artifacts {
		set.getGavOf(artifact).add(artifact)
	}

	return set.Data
}
