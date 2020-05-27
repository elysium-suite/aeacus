package main

func launchIDPrompt() {
    powerShellPrompt := `
        Add-Type -AssemblyName System.Windows.Forms
        [System.Windows.Forms.Application]::EnableVisualStyles()

        $Form                            = New-Object system.Windows.Forms.Form
        $Form.ClientSize                 = '300,162'
        $Form.text                       = "Aeacus"
        $Form.TopMost                    = $false

        $Label1                          = New-Object system.Windows.Forms.Label
        $Label1.text                     = "Enter your ID"
        $Label1.AutoSize                 = $true
        $Label1.width                    = 25
        $Label1.height                   = 10
        $Label1.location                 = New-Object System.Drawing.Point(96,30)
        $Label1.Font                     = 'Microsoft Sans Serif,10'

        $TextBox1                        = New-Object system.Windows.Forms.TextBox
        $TextBox1.multiline              = $false
        $TextBox1.width                  = 100
        $TextBox1.height                 = 20
        $TextBox1.location               = New-Object System.Drawing.Point(98,59)
        $TextBox1.Font                   = 'Microsoft Sans Serif,10'

        $Button1                         = New-Object system.Windows.Forms.Button
        $Button1.text                    = "Submit"
        $Button1.width                   = 70
        $Button1.height                  = 30
        $Button1.location                = New-Object System.Drawing.Point(111,95)
        $Button1.Font                    = 'Microsoft Sans Serif,10'

        $Form.controls.AddRange(@($Label1,$TextBox1,$Button1))
        $Button1.Add_Click({ setID })

        function setID {
            $global:id=$TextBox1.Text
            echo $id > C:\aeacus\misc\TeamID.txt
            $form.Close()
        }

        [void]$Form.ShowDialog()
    `
    shellCommand(powerShellPrompt)
}
