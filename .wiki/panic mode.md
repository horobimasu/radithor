# panic mode
panic mode attempts to make it significantly harder to remove or find radithor by doing the following:
- bsod when it detects the opening of the control center window
- disable the network selection ui on the lockscreen (requires admin)
- bsod when it detects task manager, task scheduler, registry editor, windows security and control panel
- disbale the above apps from being opened (requires admin)
- block sites that may help with remediation

# notes
1. on a device restart raditor will not trigger a bsod if the control center window is open until it is connected
