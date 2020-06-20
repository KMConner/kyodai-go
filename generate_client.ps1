Push-Location $PSScriptRoot

if (![System.IO.File]::Exists(".\openApi\openapi-generator-cli.jar")) {
    Invoke-WebRequest -OutFile openapi-generator-cli.jar `
        https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/4.3.1/openapi-generator-cli-4.3.1.jar
}

java -jar .\openApi\openapi-generator-cli.jar generate -g go `
    -i .\openApi\panda.yaml -o .\panda\client -c .\openApi\config.json --minimal-update

Pop-Location
