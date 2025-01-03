# AionApi
Aion is a habit manager that organizes, tracks, and analyzes your routine to enhance your physical and mental well-being.


- You will need to create a Base64 key for user authentication. To do this, use the function below in the main.go file and paste the result into the `.env` file under the `SECRET_KEY` variable:

```bash
func init() {
	key := make([]byte, 64)

	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}

	stringBase64 := base64.StdEncoding.EncodeToString(key)

	fmt.Println(stringBase64)
}
```

```bash
cp .env.example .env
```

