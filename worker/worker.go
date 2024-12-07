package worker

import (
	"context"
	"errors"
	"it.uniroma2.dicii/goexercise/rpc/roles"
	"sync"
)

type workerService struct {
	roles.UnimplementedRoleServiceServer
}

type role struct {
	role roles.Role
	mu   sync.Mutex
}

var (
	assignedRole = role{roles.Role_ROLE_UNKNOWN, sync.Mutex{}}
)

func (w *workerService) AssignRole(_ context.Context, assignment *roles.RoleAssignment) (*roles.RoleAssignmentResponse, error) {
	assignedRole.mu.Lock()
	defer assignedRole.mu.Unlock()
	if assignedRole.role == roles.Role_ROLE_UNKNOWN {
		// Worker role has not been assigned yet
		assignedRole.role = assignment.Role
		return &roles.RoleAssignmentResponse{Message: "Role assigned correctly"}, nil
	} else {
		return &roles.RoleAssignmentResponse{Message: "Worker role has already been assigned"}, errors.New("unable to assign worker role")
	}
}
