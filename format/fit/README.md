### Generated files
- format/fit/mappers/messages_generated.go
- format/fit/mappers/types_generated.go

### How to generate them if needed

1. Download the Fit SDK from: https://developer.garmin.com/fit/download/
2. Install NodeJS and NPM
3. Go to the `format/fit/testdata/generator` folder.
4. Run `npm install` if it's your first time 
5. Run `node index.js t /PathToSDK/Profile.xlsx > ../../mappers/types_generated.go`
6. Run `node index.js m /PathToSDK/Profile.xlsx > ../../mappers/messages_generated.go`
7. Edit `messages_generated.go` and remove the incorrect "Scale" from line ~461
8. Correct spelling of farenheit->fahrenheit and bondary->boundary to please Go linter
