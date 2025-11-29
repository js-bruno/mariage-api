package adapter_test

// func TestCreateGift(t *testing.T) {
// 	client, closeConn := adapter.MoggoOpenConnetion()()
// 	defer closeConn()
//
// 	giftNameExpected := "Presente teste"
// 	ctx := context.Background()
// 	result, err := repository.InsertGift(ctx, client)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	gift, err := repository.GetGift(ctx, client, result)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	if gift.Name != giftNameExpected {
// 		t.Errorf("%s its diferent from %s", gift.Name, giftNameExpected)
// 	}
// }
