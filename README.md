# cyanic
A tool for simple blue-green deployment

# Run
```
sudo go run . deploy-staging ./examples/dummy.yaml
```

Check command:
```
repeat 9999 (curl -s -o /dev/null -w "%{http_code}\n" localhost; sleep 0.5;)
```
