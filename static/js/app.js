(() => {
    'use strict'

    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    const forms = document.querySelectorAll('.needs-validation')

    // Loop over them and prevent submission
    Array.from(forms).forEach(form => {
      form.addEventListener('submit', event => {
        if (!form.checkValidity()) {
          event.preventDefault()
          event.stopPropagation()
        }

        form.classList.add('was-validated')
      }, false)
    })
  })();

          
  function Prompt() {
    let toastNotification = function(c) {

      const {
        msg = "", 
        icon = "success",
        pos = "top-end", 
        toast = true,
      } = c
      
      console.log("Hi love!");
      
      const Toast = Swal.mixin({
        toast: toast,
        position: pos,
        showConfirmButton: false, 
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.onmouseenter = Swal.stopTimer;
          toast.onmouseleave = Swal.resumeTimer;
        }
      });

      Toast.fire({
        icon: icon, 
        title: msg
      });
    };

    let success = function(c) {
      const {
        msg = "", 
        title = "",
        footer = ""
      } = c

      Swal.fire({
                icon: "success",
                title: title,
                text: msg,
                footer: footer
              });
    }

    let error = function(c) {
      
      const {
        msg = "", 
        title = "",
        footer = ""
      } = c

      Swal.fire({
                icon: "error",
                title: title,
                text: msg,
                footer: footer
              });
    }

    async function custom(c) {
      const {
        icon = "",
        msg = "", 
        title = "",
        showConfirmButton = true,
      } = c

      const { value: result } = await Swal.fire({
        
        icon: icon,
        title: title,
        html: msg,
        backdrop: false,
        showCancelButton: true,
        focusConfirm: false,
        showConfirmButton: showConfirmButton,

        didOpen: () => {
          if (c.didOpen !== undefined) {
            c.didOpen()
          }; 
        },
        preConfirm: () => {

          return [
          document.getElementById("start").value,
          document.getElementById("end").value
          ];
        }
        });
        
        if (result) {

        if (result.dismiss !== Swal.DismissReason.cancel) {
          if (result.value !== "") {
            if (c.callback !== undefined) {
              c.callback(result);
            }
          } else {
            c.callback(false);
          }
        } else {
          c.callback(false);
        }
        }
    }

    return {toast: toastNotification, 
      success: success, 
      error: error, 
      custom: custom}
  }



  function notify(msg, msgType) {
    notie.alert({
      type: msgType, 
      text: msg
    })
  }

  let attention = Prompt()
