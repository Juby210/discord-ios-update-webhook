app="DiUW"
files="config.example.json"

echo "Builiding $app"
go build -o build/$app
echo "Builiding $app.exe"
GOOS=windows go build -o build/$app.exe
echo "Builiding "$app"_freebsd"
GOOS=freebsd go build -o build/"$app"_freebsd

mkdir -p release
echo "Creating "$app"_linux_amd64.zip"
cp build/$app ./
zip -r release/"$app"_linux_amd64.zip $app $files
rm $app
echo "Creating "$app"_windows_amd64.zip"
cp build/$app.exe ./
zip -r release/"$app"_windows_amd64.zip $app.exe $files
rm $app.exe
echo "Creating "$app"_freebsd_amd64.zip"
cp build/"$app"_freebsd ./
zip -r release/"$app"_freebsd_amd64.zip "$app"_freebsd $files
rm "$app"_freebsd
