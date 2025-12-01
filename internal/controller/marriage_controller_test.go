package controller_test

import (
	"os"
	"slices"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/js-bruno/mariage-api/internal/adapter"
)

func TestMain(m *testing.M) {
	// pool, err := dockertest.NewPool("")
	// if err != nil {
	// 	log.Fatalf("Some shitt happens in pool connection : %s", err)
	// }
	// resource, err := pool.RunWithOptions(&dockertest.RunOptions{
	// 	Repository: "mongo",
	// 	Tag:        "latest",
	// })
	// if err != nil {
	// 	log.Fatalf("container creation : %s", err)
	// }
	// defer resource.Close()
	//
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestInsertSqlite(t *testing.T) {
	liteClient, close := adapter.NewSqliteClient("file:%s?mode=memory&cache=shared")
	defer close()

	gift := adapter.Gift{
		ID:        1,
		Name:      "Microondas Philco PMO30EP 28 Litros Preto Com Porta Espelhada 1400w",
		Price:     649.90,
		Reserved:  false,
		Category:  "Cozinha",
		Buyers:    0,
		MaxBuyers: 1,
		Image:     "https://a-static.mlcdn.com.br/420x420/microondas-philco-pmo30ep-28-litros-preto-com-porta-espelhada-1400w/techshop/frnphc00095/657e6a8da6b8501acf203d7b8d731159.jpeg",
		Link:      "https://www.magazineluiza.com.br/microondas-philco-pmo30ep-28-litros-preto-com-porta-espelhada-1400w/p/aecfake718/ed/mond/?seller_id=techshop",
		QRCode:    "",
	}
	err := liteClient.InsertGift(gift)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSelectSqlite(t *testing.T) {
	giftNames := []string{"Microondas", "Lancheira", "Crocodilo", "Computador"}
	liteClient, close := adapter.NewSqliteClient("file:%s?mode=memory&cache=shared")
	defer close()
	preInsertGifts(liteClient, giftNames)

	gifts, err := liteClient.ListGifts()
	if err != nil {
		t.Error(err)
	}
	for _, g := range gifts {
		if !slices.Contains(giftNames, g.Name) {
			t.Errorf("%s does not exists in the list", g.Name)
		}

	}

}

func preInsertGifts(client *adapter.SqliteClient, names []string) (*[]adapter.Gift, error) {
	gifts := []adapter.Gift{}
	for i, name := range names {
		gift := adapter.Gift{
			ID:        i,
			Name:      name,
			Price:     649.90,
			Reserved:  false,
			Category:  "Cozinha",
			Buyers:    0,
			MaxBuyers: 1,
			Image:     "https://a-static.mlcdn.com.br/420x420/microondas-philco-pmo30ep-28-litros-preto-com-porta-espelhada-1400w/techshop/frnphc00095/657e6a8da6b8501acf203d7b8d731159.jpeg",
			Link:      "https://www.magazineluiza.com.br/microondas-philco-pmo30ep-28-litros-preto-com-porta-espelhada-1400w/p/aecfake718/ed/mond/?seller_id=techshop",
			QRCode:    "12312312i8",
		}
		err := client.InsertGift(gift)
		if err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}

	return &gifts, nil
}
