package globalref

import (
	"github.com/hashicorp/terraform/internal/addrs"
)

// ReferencesFromResource returns all of the direct references from the
// definition of the resource instance at the given address. It doesn't
// include any indirect references.
//
// Resource configurations can only refer to other objects within the same
// module, so callers should assume that the returned references are all
// relative to the same module instance that the given address belongs to.
func (a *Analyzer) ReferencesFromResourceInstance(addr addrs.AbsResourceInstance) []*addrs.Reference {
	// Using MetaReferences for this is kinda overkill, since
	// lang.ReferencesInBlock would be sufficient really, but
	// this ensures we keep consistent and aside from some
	// extra overhead this call boils down to a call to
	// lang.ReferencesInBlock anyway.
	fakeRef := &addrs.Reference{
		Subject: addr.Resource,
	}
	_, refs := a.MetaReferences(addr.Module, fakeRef)
	return refs
}
