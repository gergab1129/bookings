function SearchAvailability(CSRFToken, RoomId) {


        let book_btn = document.getElementById("check-availability-btn")
    
        book_btn.addEventListener("click", function() {
            let html = `<form id="check-avilability-form" action="" method="post" novalidate class="needs-validation">
                        <div class="row">
                            <div class="col">
                                <div id="date-picker">
                                <div class="row">
                                    <div class="col">
                                    <input required type="text" class="form-control" name="start" id="start" placeholder="Arrival Date">
                                    </div>
                                    <div class="col">
                                    <input required type="text" class="form-control" name="end" id="end" placeholder="Departure Date">  
                                    </div>
                                </div>
                                </div>
                            </div>
                        </div>
                    </form>`
    
            attention.custom({
            title: "Select date range"
            , msg: html
            , callback: function(result) {
                console.log("called")
    
                let form = document.getElementById("check-avilability-form")
                let formData = new FormData(form)
                formData.append("csrf_token", CSRFToken);
                formData.append("room_id", RoomId);
    
                fetch('/search-availability-json', {
                    method: "post", 
                    body: formData
                }) 
                .then(response => response.json())
                .then(data => {
                    if (data.ok) {
                        attention.custom({
                            icon: "success",
                            msg: '<p>Room is available</p>'
                                + '<p><a href="/book-room?id='
                                 + data.room_id 
                                 + '&s='
                                 + data.start_date 
                                 + '&e='
                                 + data.end_date
                                  + '" class="btn btn-primary">'
                                 + 'Book now</a></p>',
                            showConfirmButton: false,
    
                        })
                    } else {
                        attention.error({
                            msg: "No Availability",
                        })
                    }
                })
    
            }
        ,  didOpen: () => {
            const datePicker = document.getElementById("date-picker")
            const rp = new DateRangePicker(datePicker, {
            buttonClass: 'btn',
            format: 'yyyy-mm-dd',
            showOnFocus: true,
            orientation: "top", 
            showOnClick: true, 
            minDate: new Date()
    
            // ...options
          }); 
          }
        })
            
        })   
  }