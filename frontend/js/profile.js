//Profile JS file Is emty now
const ActiveUser = document.getElementById("UserName")
const ProfileUName = document.getElementById("UserInfoName")
const UserEmailInfo = document.getElementById("UserEmailInfo")
const UserRole = document.getElementById("UserRole")

export function checkUserEmail(){
    let PName = ProfileUName.innerHTML
    let AName = ActiveUser.innerHTML
    let Role = UserRole.innerHTML
    if (Role != "Admin"){
        if (PName != AName ){
            UserEmailInfo.remove()
        }
    }
}

