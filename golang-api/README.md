
```
# 1. build da imagem Docker
docker build -t app-desafio-go .

# 2. Inicie um novo contêiner com a imagem reconstruída
docker run -d -p 8080:8080  app-desafio-go

```