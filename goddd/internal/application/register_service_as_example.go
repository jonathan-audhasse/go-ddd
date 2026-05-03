package application

// import (
// 	"context"
// 	"yourapp/domain"
// )

// type RegisterUserService struct {
// 	tx       TransactionManager
// 	users    domain.UserRepository
// 	accounts domain.AccountRepository
// }

// func NewRegisterUserService(
// 	tx TransactionManager,
// 	users domain.UserRepository,
// 	accounts domain.AccountRepository,
// ) *RegisterUserService {
// 	return &RegisterUserService{
// 		tx:       tx,
// 		users:    users,
// 		accounts: accounts,
// 	}
// }

// func (s *RegisterUserService) Register(ctx context.Context, name string) error {
// 	return s.tx.Do(ctx, func(ctx context.Context) error {
// 		user := domain.NewUser(name)

// 		if err := s.users.Save(ctx, user); err != nil {
// 			return err
// 		}

// 		if err := s.accounts.CreateForUser(ctx, user.ID); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }
