# Automated testing

## 1. Initialize module `validator`
```
go mod init example.com/validator
```

## 2. Get testify lib
```
go get github.com/stretchr/testify
```

## 3. Add missing libs in go.sum
```
go mod tidy
```

## 4. Test
```
go test -v
```

## Notes
```
# Download modules to local cache if you did not get libs
go mod download

# Then run test
go test -v
```