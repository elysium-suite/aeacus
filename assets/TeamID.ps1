# Powershell script to prompt user to enter TeamID

$teamIDContent = Get-Content C:\Users\tanay\Programing\Powershell\TeamID.txt

if ($teamIDContent -eq "YOUR-TEAMID-HERE") {
    Start-Process -FilePath C:\aeacus\aeacus.exe -ArgumentList "idprompt"
}