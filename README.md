# az-dns-update

This is a simple script to update a DNS record in Azure DNS with the current public ip of the machine.

Environment:
- AZURE_CLIENT_ID
- AZURE_CLIENT_SECRET
- AZURE_TENANT_ID
- AZURE_SUBSCRIPTION_ID
- KEEP_ALIVE

For sub-domain:
```
docker rm -f dynamic-dns ; \
docker run -d --restart always \
  --env AZURE_CLIENT_ID=xxxxxxxx \
  --env AZURE_CLIENT_SECRET=xxxxxxxx \
  --env AZURE_TENANT_ID=xxxxxxxx \
  --env AZURE_SUBSCRIPTION_ID=xxxxxxxx \
  --env KEEP_ALIVE=true \
  --name dynamic-dns \
  foilen/foilen-cloud-tools:latest /usr/bin/az-dns-update \
    myResourceGroup example.com test.example.com && \
docker logs -f dynamic-dns
```

For main domain:
```
docker rm -f dynamic-dns ; \
docker run -d --restart always \
  --env AZURE_CLIENT_ID=xxxxxxxx \
  --env AZURE_CLIENT_SECRET=xxxxxxxx \
  --env AZURE_TENANT_ID=xxxxxxxx \
  --env AZURE_SUBSCRIPTION_ID=xxxxxxxx \
  --env KEEP_ALIVE=true \
  --name dynamic-dns \
  foilen/foilen-cloud-tools:latest /usr/bin/az-dns-update \
    myResourceGroup example.com @.example.com && \
docker logs -f dynamic-dns
```
