$Alice = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Node_Alice";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    go run . :9000 9001 Alice;
}'

invoke-expression -Command $Alice

$Bob = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Node_Bob";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    go run . :9001 9000 Bob;
}'

invoke-expression -Command $Bob