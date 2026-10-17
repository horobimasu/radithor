# panic mode
panic mode attempts to make it significantly harder to remove or find radithor by doing the following:
- bsod when it detects the opening of the control center window
- disable the network selection ui on the lockscreen (requires admin)
- bsod when it detects task manager, task scheduler, registry editor, windows security and control panel
- disbale the above apps from being opened (requires admin)
- block sites that may help with remediation
- bsod when it detects processes titles with certain words that may suggest the user is trying to remove radithor

# notes
1. on a device restart raditor will not trigger a bsod if the control center window is open until it is connected
2. when it detects task manager, task scheduler, the certain words in the title etc, it will attempt to kill them first, if that does not work it will trigger a bsod
