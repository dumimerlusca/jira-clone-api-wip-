package consts

const (
	// AUTH
	ApiPathLogin    = "/api/auth/login"
	ApiPathRegister = "/api/auth/register"

	// TICKETS
	ApiPathCreateTicker = "/api/tickets"
	ApiPathUpdateTicket = "/api/tockets/{id}"

	// PROJECTS
	ApiPathCreateProject     = "/api/projects/create"
	ApiPathUpdateProject     = "/api/projects/update/{projectId}"
	ApiPathGetProjectDetails = "/api/projects/details/{projectId}"

	// PROJECT INVITATIONS
	ApiPathSendProjectInvite         = "/api/projects/sendInvite/{projectId}"
	ApiPathAcceptProjectInvite       = "/api/projects/acceptInvite/{inviteId}"
	ApiPathRejectProjectInvite       = "/api/projects/rejectInvite/{inviteId}"
	ApiPathGetSentProjectInvites     = "/api/project-invites/sent"
	ApiPathGetReceivedProjectInvites = "/api/project-invites/received"
)
