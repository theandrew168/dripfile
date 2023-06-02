package cli

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/theandrew168/dripfile/internal/location"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Application struct {
	args []string

	locationService location.Service
	transferService transfer.Service
}

func New(args []string, locationService location.Service, transferService transfer.Service) *Application {
	app := Application{
		args: args,

		locationService: locationService,
		transferService: transferService,
	}
	return &app
}

func (app *Application) Run(ctx context.Context) error {
	err := app.dripfile(app.args)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (app *Application) dripfile(args []string) error {
	if len(args) == 0 {
		fmt.Println("usage: dripfile [location transfer]")
		return nil
	}

	cmd := args[0]
	switch cmd {
	case "location":
		return app.location(args[1:])
	case "transfer":
		return app.transfer(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (app *Application) location(args []string) error {
	if len(args) == 0 {
		fmt.Println("usage: dripfile location [get add remove]")
		return nil
	}

	cmd := args[0]
	switch cmd {
	case "get":
		return app.locationGet(args[1:])
	case "add":
		return app.locationAdd(args[1:])
	case "remove":
		return app.locationRemove(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (app *Application) locationGet(args []string) error {
	var ls []*location.Location
	if len(args) == 0 {
		var err error
		ls, err = app.locationService.GetAll(location.GetAllQuery{})
		if err != nil {
			return err
		}
	} else {
		id := args[0]
		l, err := app.locationService.GetByID(location.GetByIDQuery{
			ID: id,
		})
		if err != nil {
			return err
		}

		ls = append(ls, l)
	}

	for _, l := range ls {
		switch l.Kind() {
		case location.KindMemory:
			fmt.Printf(
				"%s %s\n",
				l.ID(),
				l.Kind(),
			)
		case location.KindS3:
			info := l.S3Info()
			fmt.Printf(
				"%s %s %s %s\n",
				l.ID(),
				l.Kind(),
				info.Endpoint,
				info.Bucket,
			)
		}
	}

	return nil
}

func (app *Application) locationAdd(args []string) error {
	endpoint, err := input("Endpoint: ")
	if err != nil {
		return err
	}
	bucket, err := input("Bucket: ")
	if err != nil {
		return err
	}
	accessKeyID, err := input("AccessKeyID: ")
	if err != nil {
		return err
	}
	secretAccessKey, err := input("SecretAccessKey: ")
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	err = app.locationService.AddS3(location.AddS3Command{
		ID: id.String(),

		Endpoint:        endpoint,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	})
	if err != nil {
		return err
	}

	fmt.Printf("location created: %s\n", id)
	return nil
}

func (app *Application) locationRemove(args []string) error {
	if len(args) == 0 {
		fmt.Println("usage: dripfile location remove [id]")
		return nil
	}

	id := args[0]
	return app.locationService.Remove(location.RemoveCommand{
		ID: id,
	})
}

func (app *Application) transfer(args []string) error {
	if len(args) == 0 {
		fmt.Println("usage: dripfile transfer [get add remove run]")
		return nil
	}

	cmd := args[0]
	switch cmd {
	case "get":
		return app.transferGet(args[1:])
	case "add":
		return app.transferAdd(args[1:])
	case "remove":
		return app.transferRemove(args[1:])
	case "run":
		return app.transferRun(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (app *Application) transferGet(args []string) error {
	var ts []*transfer.Transfer
	if len(args) == 0 {
		var err error
		ts, err = app.transferService.GetAll(transfer.GetAllQuery{})
		if err != nil {
			return err
		}
	} else {
		id := args[0]
		t, err := app.transferService.GetByID(transfer.GetByIDQuery{
			ID: id,
		})
		if err != nil {
			return err
		}

		ts = append(ts, t)
	}

	for _, t := range ts {
		fmt.Printf(
			"%s %s\n",
			t.ID(),
			t.Pattern(),
		)
	}

	return nil
}

func (app *Application) transferAdd(args []string) error {
	pattern, err := input("Pattern: ")
	if err != nil {
		return err
	}
	fromLocationID, err := input("FromLocationID: ")
	if err != nil {
		return err
	}
	toLocationID, err := input("ToLocationID: ")
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	err = app.transferService.Add(transfer.AddCommand{
		ID: id.String(),

		Pattern:        pattern,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("transfer created: %s\n", id)
	return nil
}

func (app *Application) transferRemove(args []string) error {
	if len(args) == 0 {
		fmt.Println("usage: dripfile transfer remove [id]")
		return nil
	}

	id := args[0]
	return app.transferService.Remove(transfer.RemoveCommand{
		ID: id,
	})
}

func (app *Application) transferRun(args []string) error {
	if len(args) == 0 {
		fmt.Println("usage: dripfile transfer run [id]")
		return nil
	}

	id := args[0]
	return app.transferService.Run(transfer.RunCommand{
		ID: id,
	})
}

func input(prompt string) (string, error) {
	fmt.Print(prompt)

	var resp string
	_, err := fmt.Scanln(&resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}
