package router

import (
	"github.com/chinathaip/chatify/service"
	"github.com/google/uuid"
)

type mockChatService struct {
	isGetAllCalled bool
}

func (cs *mockChatService) GetAllChat() ([]service.Chat, error) {
	cs.isGetAllCalled = true
	return []service.Chat{
		{
			ID:   1,
			Name: "Test Chat Room",
		},
		{
			ID:   2,
			Name: "Test Chat Room 2",
		},
	}, nil
}

func (cs *mockChatService) DeleteChat(chatID int) error {
	return nil
}

func (cs *mockChatService) CreateNewChat(*service.Chat) error {
	return nil
}

func (cs *mockChatService) IsChatExist(chatName string) (int, bool) {
	return 0, true
}

type mockMessageService struct {
	isGetMessagesCalled     bool
	isStoreNewMessageCalled bool
}

type mockError struct{}

func (e *mockError) Error() string {
	return "Error occured!"
}

func (ms *mockMessageService) GetMessagesInChat(chatID, pageNumber, pageSize int) ([]service.Message, error) {
	ms.isGetMessagesCalled = true
	if chatID == 1 {
		return []service.Message{
			{
				ID:     1,
				ChatID: 1,
				Data:   "Message 1",
			},
			{
				ID:     2,
				ChatID: 1,
				Data:   "Message 2 from the same dude",
			},
		}, nil
	}
	return []service.Message{}, &mockError{}
}

func (ms *mockMessageService) StoreNewMessage(msg *service.Message) error {
	ms.isStoreNewMessageCalled = true
	return nil
}

type mockUserService struct {
	isGetUserNameByIDCalled bool
	isCreateNewUserCalled   bool
}

func (ms *mockUserService) GetUserNameByID(id uuid.UUID) (string, error) {
	ms.isGetUserNameByIDCalled = true
	return "", nil
}

func (ms *mockUserService) CreateNewUser(user service.User) error {
	ms.isCreateNewUserCalled = true
	return nil
}
