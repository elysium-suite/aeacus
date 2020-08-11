package main

func launchIDPrompt() {
	powerShellPrompt := `
    Add-Type -AssemblyName System.Windows.Forms
    [System.Windows.Forms.Application]::EnableVisualStyles()
    
    $Form                            = New-Object system.Windows.Forms.Form
    $Form.ClientSize                 = New-Object System.Drawing.Point(520,265)
    $Form.text                       = "Aeacus"
    $Form.TopMost                    = $true
    $Form.Icon                       = "C:\aeacus\assets\logo.ico"
    $Form.BackgroundImage            = [system.drawing.image]::FromFile("C:\aeacus\assets\background_520x265.png")
    
    $Label1                          = New-Object system.Windows.Forms.Label
    $Label1.text                     = "Enter Your Unique Team ID"
    $Label1.AutoSize                 = $true
    $Label1.width                    = 170
    $Label1.height                   = 9
    $Label1.location                 = New-Object System.Drawing.Point(117,44)
    $Label1.Font                     = New-Object System.Drawing.Font('Raleway',16,[System.Drawing.FontStyle]([System.Drawing.FontStyle]::Bold))
    $Label1.BackColor                = [System.Drawing.Color]::FromName("Transparent")
    $Label1.ForeColor                = [System.Drawing.ColorTranslator]::FromHtml("#ffffff")
    
    $TextBox1                        = New-Object system.Windows.Forms.TextBox
    $TextBox1.multiline              = $false
    $TextBox1.width                  = 333
    $TextBox1.height                 = 38
    $TextBox1.location               = New-Object System.Drawing.Point(80,109)
    $TextBox1.Font                   = New-Object System.Drawing.Font('Consolas',12)
    
    $Button1                         = New-Object system.Windows.Forms.Button
    $Button1.text                    = "Validate"
    $Button1.width                   = 110
    $Button1.height                  = 40
    $Button1.location                = New-Object System.Drawing.Point(202,155)
    $Button1.Font                    = New-Object System.Drawing.Font('Raleway',10)
    $Button1.Image                   = [System.Drawing.Image]::FromFile("C:\aeacus\assets\pp_109x36.png")
    $Button1.ForeColor               = [System.Drawing.ColorTranslator]::FromHtml("#ffffff")
    
    $Form.controls.AddRange(@($Label1,$TextBox1,$Button1))
    $Button1.Add_Click({ setID })

    function setID {
        $global:id=$TextBox1.Text
        echo $id > C:\aeacus\TeamID.txt
        $form.Close()
    }

    [void]$Form.ShowDialog()
    `
	shellCommand(powerShellPrompt)
}

func launchConfigGui() {
	warnPrint("This feature is not supported yet on Windows.")
}
