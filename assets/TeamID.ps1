# Powershell script to prompt user to enter TeamID

$teamIDContent = Get-Content C:\aeacus\TeamID.txt

if ($teamIDContent -eq "YOUR-TEAMID-HERE") {
    Start-Process -FilePath C:\aeacus\phocus.exe -ArgumentList "idprompt"
}