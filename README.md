# cyanic
A tool for simple blue-green deployment

# Run
Using the default config location:
```
sudo go run . deploy-staging
```

Using custom config location:
```
sudo go run . deploy-staging ./path/to/cyanic/config.yaml
```

# Misc
Check command port:
```
repeat 9999 (curl -s -o /dev/null -w "%{http_code}\n" localhost:80; sleep 0.5;)
```
