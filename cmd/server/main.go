package main

func main() {
	//if err := godotenv.Load(".env.local"); err != nil {
	//	log.Printf("Failed to load .env.local.\n")
	//}
	//dbPath := os.Getenv("DB_DSN")
	//if dbPath == "" {
	//	log.Fatal("DB_DSN environment variable not set!")
	//}
	//fmt.Println(" Connecting with DSN:", dbPath)
	//db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	//if err != nil {
	//	log.Fatal("❌ failed to connect database: ", err)
	//}
	//
	//m := backup.NewMigrator(db)
	//m.AddMigration(backup.CreateOrganizationTable())
	//m.AddMigration(backup.CreateLoanApplicationTable())
	//
	//if len(os.Args) < 2 {
	//	log.Fatal("Usage: go run main.go [up|down|status]")
	//}
	//
	//cmd := os.Args[1]
	//
	//switch cmd {
	//case "up":
	//	if err := m.Up(); err != nil {
	//		log.Fatal("❌ Migration UP failed: ", err)
	//	}
	//	fmt.Println("✅ All migrations applied successfully!")
	//
	//case "down":
	//	if err := m.Down(); err != nil {
	//		log.Fatal("❌ Migration DOWN failed: ", err)
	//	}
	//	fmt.Println("✅ Last migration rolled back successfully!")
	//
	//case "status":
	//	if err := m.Status(context.Background()); err != nil {
	//		log.Fatal("❌ Status check failed: ", err)
	//	}
	//
	//default:
	//	log.Fatal("Unknown command: ", cmd)
	//}
	//DB_DSN="host=localhost user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} port=${POSTGRES_PORT} sslmode=${POSTGRES_SSL_MODE}"

}
