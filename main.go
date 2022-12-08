package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"entdemo/ent"
	"entdemo/ent/car"
	"entdemo/ent/user"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client, err := ent.Open("mysql", "yuki:#Yuk1@23@tcp(localhost:3306)/sample_golang?parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	//if err := QueryCars(ctx, a8m); err != nil {
	//	log.Fatal(err)
	//}

	router := gin.Default()

	router.GET("/test-router", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	router.POST("/create-user", func(c *gin.Context) {
		name := c.PostForm("name")
		if _, err = CreateUser(ctx, client, name); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		}
	})
	router.GET("/query-user", func(c *gin.Context) {
		name := c.Query("name")
		u, err := QueryUser(ctx, client, name)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": u,
		})
	})
	router.POST("/create-car", func(c *gin.Context) {
		name := c.PostForm("name")
		a8m, err := CreateCars(ctx, client, name)
		if err != nil {
			if err != nil {
				log.Fatal(err)
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": a8m,
			})
		}
	})
	router.GET("/list-car", func(c *gin.Context) {
		name := c.Query("name")
		cars, err := QueryCars(ctx, client, name)
		if err != nil {
			if err != nil {
				log.Fatal(err)
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": cars,
			})
		}
	})

	//router.PUT("/somePut", putting)
	//router.DELETE("/someDelete", deleting)
	//router.PATCH("/somePatch", patching)
	//router.HEAD("/someHead", head)
	//router.OPTIONS("/someOptions", options)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

func CreateUser(ctx context.Context, client *ent.Client, name string) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client, name string) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ(name)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client, name string) (*ent.User, error) {
	// Create a new car with model "Tesla".
	tesla, err := client.Car.
		Create().
		SetModel(name).
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("Ta8ms").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, client *ent.Client, name string) ([]*ent.Car, error) {
	cars, err := client.Car.Query().All(ctx)
	user, _ := QueryUser(ctx, client, name)
	if err != nil {
		return nil, fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars:", cars)

	// What about filtering specific cars.
	ford, err := user.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println(ford)
	return cars, nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}

	// Query the inverse edge.
	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
		}
		log.Printf("car %q owner: %q\n", c.Model, owner.Name)
	}
	return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	// First, create the users.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("Ariel").
		Save(ctx)
	if err != nil {
		return err
	}
	neta, err := client.User.
		Create().
		SetAge(28).
		SetName("Neta").
		Save(ctx)
	if err != nil {
		return err
	}
	// Then, create the cars, and attach them to the users created above.
	err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		// Attach this car to Ariel.
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		// Attach this car to Ariel.
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		// Attach this graph to Neta.
		SetOwner(neta).
		Exec(ctx)
	if err != nil {
		return err
	}
	// Create the groups, and add their users in the creation.
	err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(neta, a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("The graph was created successfully")
	return nil
}
