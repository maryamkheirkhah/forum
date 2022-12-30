const pwShow = document.querySelector(".show");
userName = document.querySelector("#username")
userPassword = document.querySelector("#user_password")
userConfirmPassword = document.querySelector("#confirm_password")
alertUsername = document.querySelector(".alert_Username")
alertUsernameIcon = document.querySelector(".alertUIcon")
alertEmail = document.querySelector(".alert_Email")
alertEmailIcon = document.querySelector(".alertEmailIcon")
alertPass = document.querySelector(".alert_Password")
alertPassIcon = document.querySelector(".alertPassIcon")

let checkUsername = false
let checkConfirmPassword = false
let checkEmail = false

// js code to check #username length
userName.addEventListener("input",()=>{
    let inputValue = userName.value.trim();// trim function
    if((inputValue.length < 8)){
        checkUsername = false
        alertUsername.innerText = "Username enter at least 8 characters"
        alertUsername.style.color ="red"
        alertUsernameIcon.style.display = "inline"
        alertUsernameIcon.style.color = "red"
        submitButton.setAttribute("disabled",true);
        document.getElementById("alertU").className ="fas fa-exclamation-circle alertUIcon"
       
    }else if((inputValue.length >=8)==((/^[A-Za-z][A-Za-z0-9_]{8,24}$/.test(userName.value)))){
        checkUsername = true
        alertUsername.innerText = "Username format correct"
        alertUsername.style.color ="green"
        alertUsernameIcon.style.display = "inline"
        alertUsernameIcon.style.color = "green"
        document.getElementById("alertU").className = "fa-solid fa-circle-check alertUIcon"
                //if password and email format correct open submit button
                if ((checkEmail != false)&&(checkEmail==checkConfirmPassword)&&(checkUsername)){    
                    submitButton.removeAttribute("disabled")
                }else{
                    submitButton.setAttribute("disabled",true);
                    console.log(checkUsername,"user",checkEmail,"email",checkConfirmPassword,"pass")
                }
                //check specialCharacter
    }else if (!(/^[A-Za-z][A-Za-z0-9]{7,24}$/.test(userName.value))) {
        checkUsername = false
        alertUsername.innerText = "Username should not contain special characters"
        alertUsername.style.color ="red"
        alertUsernameIcon.style.display = "inline"
        alertUsernameIcon.style.color = "red"
        submitButton.setAttribute("disabled",true);
        document.getElementById("alertU").className ="fas fa-exclamation-circle alertUIcon"
    }

})
// js code to check #email format
function validateEmail(emailId){
var mailformat = /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/;
if(emailId.value.match(mailformat))
{
    checkEmail = true
    alertEmail.innerText = "Email format correct"
    alertEmail.style.color = "green"
    alertEmailIcon.style.display = "inline"
    alertEmailIcon.style.color = "green"
    document.getElementById("alertE").className = "fa-solid fa-circle-check alertEmailIcon"
    //if password and email format correct open submit button
    if ((checkEmail != false)&&(checkEmail==checkConfirmPassword)&&(checkUsername)){    
        submitButton.removeAttribute("disabled")
    }else{
        submitButton.setAttribute("disabled",true);
        console.log(checkUsername,"user",checkEmail,"email",checkConfirmPassword,"pass")
    }
}
else
{
    checkEmail = false
    alertEmail.innerText = "Email format not correct"
    alertEmailIcon.style.display = "inline"
    alertEmail.style.color = "red"
    alertEmailIcon.style.color = "red"
    submitButton.setAttribute("disabled",true);
    document.getElementById("alertE").className ="fas fa-exclamation-circle alertEmailIcon"
}
}


//js code to check password and confirm password is same or not
var passConfirm = function(){
    if ((document.getElementById("user_password").value == 
        document.getElementById("confirm_password").value)&&(document.getElementById("user_password").value.length!=0)&&(document.getElementById("user_password").value.length>=8)){
        alertPass.innerText = "Password matched";
        alertPass.style.color = "green";
        alertPassIcon.style.color = "green";
        document.getElementById("alertPass").className = "fa-solid fa-circle-check alertPassIcon"
        alertPassIcon.style.display = "inline";
        checkConfirmPassword=true
        //if password and email format correct open submit button
        if ((checkEmail != false)&&(checkEmail==checkConfirmPassword)&&(checkUsername)){    
            submitButton.removeAttribute("disabled")
        }else{
            submitButton.setAttribute("disabled",true);
            console.log(checkUsername,"user",checkEmail,"email",checkConfirmPassword,"pass")
        }
    }else if (document.getElementById("user_password").value.length==0){
        alertPass.innerText = "Enter  8 - 24 characters";
        alertPass.style.color = "white";
        alertPassIcon.style.display = "none";
        submitButton.setAttribute("disabled",true);
    }else{
        checkConfirmPassword=false
        alertPass.innerText = "Password did not match";
        alertPass.style.color = "red";
        alertPassIcon.style.color = "red";
        alertPassIcon.style.display = "inline";
        document.getElementById("alertPass").className ="fas fa-exclamation-circle alertPassIcon"
        submitButton.setAttribute("disabled",true);
    }
}

// js code to check and confirm input field's password
userPassword.addEventListener("input",()=>{
    let inputValue = userPassword.value.trim();// trim function
    if(inputValue.length >= 8){
        alertPass.innerText = "enter at least 8 characters";
        userConfirmPassword.removeAttribute("disabled");
    }else if ((document.getElementById("user_password").value.length==0)==(document.getElementById("confirm_password").value.length!=0)){ 
        userConfirmPassword.removeAttribute("disabled");
    }else{
        userConfirmPassword.setAttribute("disabled",true);
        //submitButton.setAttribute("disabled",true);
        checkConfirmPassword = false
        userConfirmPassword.Value = "";   
    }
});

//js code show and hide pasword
pwShow.addEventListener("click", ()=>{
    if ((userPassword.type === "password")&&(userPassword.type == "password")){
        userConfirmPassword.type = "text";
        userPassword.type = "text";
        pwShow.classList.replace("fa-eye-slash", "fa-eye")
    }else {
        userConfirmPassword.type = "password";
        userPassword.type = "password";
        pwShow.classList.replace( "fa-eye","fa-eye-slash")
    }
});





