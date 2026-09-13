package main

import "fmt"

type ModerationError struct {
	ContentId string
	Reason    string
}

func (e ModerationError) Error() string {
	return fmt.Sprintf("контент %s отклонён: %s", e.ContentId, e.Reason)
}

func CheckPoster(status string) error {
	if status == "approved" {
		return nil
	}
	return ModerationError{
		ContentId: "poster-42",
		Reason:    fmt.Sprintf("статус модерации: %s", status),
	}
}

func main() {
	err := CheckPoster("approved")
	fmt.Printf("approved: err=%v, err == nil: %v\n", err, err == nil)

	err = CheckPoster("rejected")
	fmt.Printf("rejected: %v\n", err)

	if me, ok := err.(ModerationError); ok {
		fmt.Printf("type assertion: ContentId=%s, Reason=%s\n", me.ContentId, me.Reason)
	}

	var ptr *ModerationError
	var iface error = ptr
	fmt.Printf("nil pointer in error: iface=%v, iface == nil: %v\n", iface, iface == nil)
	if iface != nil {
		fmt.Println("интерфейс не nil: тип *ModerationError, значение указателя nil")
	}
}
