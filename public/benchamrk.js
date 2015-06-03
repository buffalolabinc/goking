(function(){

    var apiBenchmark = require('api-benchmark');
     
    var service = {
      goqueenapi: "http://localhost:8080/"
    };

    var options = { 
        debug: true,
        minSamples: 1000, 
        stopOnError: false,
        maxTime: 40,
        runMode: 'parallel',
        maxConcurrentRequests: 500,
    };
     
    var routes = { 
        getSchedules: {
            method: 'get', 
            route: 'api/schedules',
        },
        getScheduleById: {
            method: 'get', 
            route: 'api/schedules/' + (function() { 
                return randomnumber=Math.floor(Math.random()*100);
            })(),
        },
        postSchedules: {
            method: 'post',
            route: 'api/schedules',
            data: getScheduleData()
        },
        stats: {
            method: 'get',
            route: 'stats' 
        }
    };
    var start = new Date().getTime();



    apiBenchmark.compare(service, routes, options, function(err, results){
      console.log(results, err);
    var end = new Date().getTime();
    var time = end - start;

    console.log(time/100)
      // displays some stats! 
    });

    function getScheduleData(){
        return {
            name:"dsfadf",
            mon:false,
            tue:false,
            wed:true,
            thu:false,
            fri:true,
            sat:true,
            sun:false,
            startTime: 1431403200,
            endTime:1431489480 
        };
    }

})()
