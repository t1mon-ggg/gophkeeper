package storage

// var (
// 	ctrl *gomock.Controller
// 	s    *storage.MockStorage
// 	t    *testing.T
// )

// func init() {
// 	t = new(testing.T)
// 	ctrl = gomock.NewController(t)
// 	s = storage.NewMockStorage(ctrl)
// 	go func() {
// 		defer ctrl.Finish()
// 	}()
// 	gomock.InOrder(
// 		s.EXPECT().Ping().Return(errors.New("connection lost")),
// 		s.EXPECT().Ping().Return(nil),
// 	)
// 	gomock.InOrder(
// 		s.EXPECT().SignUp(gomock.Any().String(), gomock.Any().String(), gomock.Any()).Return(errors.New("wrong username or password")),
// 		s.EXPECT().SignUp(gomock.Any().String(), gomock.Any().String(), gomock.Any()).Return(errors.New("db connection error")),
// 		s.EXPECT().SignUp(gomock.Any().String(), gomock.Any().String(), gomock.Any()).Return(nil),
// 	)

// 	gomock.InOrder(
// 		s.EXPECT().SaveLog(gomock.Any().String(), gomock.Any().String(), gomock.Any().String(), gomock.Any(), gomock.Any()).Return(nil),
// 		s.EXPECT().SaveLog(gomock.Any().String(), gomock.Any().String(), gomock.Any().String(), gomock.Any(), gomock.Any()).Return(nil),
// 		s.EXPECT().SaveLog(gomock.Any().String(), gomock.Any().String(), gomock.Any().String(), gomock.Any(), gomock.Any()).Return(nil),
// 	)
// }
