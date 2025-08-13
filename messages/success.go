package messages

import "fmt"

func MsgCallbackHandled(operation string) string {
	return fmt.Sprintf("%s callback handled successfully", operation)
}

func MsgRedirected(url string) string {
	return fmt.Sprintf("Redirected to %s", url)
}

func MsgCreated(entity string) string {
	return fmt.Sprintf("%s created successfully", entity)
}

func MsgFetched(entity string) string {
	return fmt.Sprintf("%s fetched successfully", entity)
}

func MsgUpdated(entity string) string {
	return fmt.Sprintf("%s updated successfully", entity)
}

func MsgDeleted(entity string) string {
	return fmt.Sprintf("%s deleted successfully", entity)
}

func MsgCancelled(entity string) string {
	return fmt.Sprintf("%s cancelled successfully", entity)
}
