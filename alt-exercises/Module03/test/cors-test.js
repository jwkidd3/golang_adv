function main() 
{
    console.log("call that require cors enabled");
    $.ajax
    ({
        dataType: "tex/plain",
        url: "http://localhost:5000/users",
        success: function(data) 
        {
            console.log("log response on success");
            console.log(data);
        }
    });
}