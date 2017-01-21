function notFound() {
    $('body').text('Not found!');
}

function handleAjaxError(jqXHr, textStatus) {
    var message = '';
 
    switch (textStatus) {
        case 'notmodified':
            message = 'Not Modified';
            break;
        case 'parsererror':
            message = 'Parser Error';
            break;
        case 'timeout':
            message = 'Time Out';
            break;
        default:
            switch (jqXHr.status) {
                case 401: // unauthorized
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '401 Unauthorized';
                    }
                    
                    break;
                case 403: // forbidden
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '403 Forbidden';
                    }
                    
                    break;
                case 404: // not found
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '404 Not Found';
                    }
                    
                    break;
                case 500: // internal server error
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '500 Internal Server Error';
                    }
                    
                    break;
                case 503: // service unavailable
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '503 Service Unavailable';
                    }
                    
                    break;
                case 555: // error
                    if (jqXHr.responseJSON) {
                        message = jqXHr.responseJSON.message;
                    }
                    else {
                        message = '555 Error';
                    }
                    
                    break;
                default:
                    message = 'Error';
            }
    }
    
    if (message) {
        console.log('Error: ' + message);
    }
}