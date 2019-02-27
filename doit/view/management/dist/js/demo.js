$(document).ready(function () {
    /**
     * Line Chart
     */


     function getQueryVariable(variable)
        {
           var query = window.location.search.substring(1);
           var vars = query.split("&");
           for (var i=0;i<vars.length;i++) {
                   var pair = vars[i].split("=");
                   if(pair[0] == variable){return pair[1];}
           }
           return(false);
        }

    var token = getQueryVariable("token")
    var year = getQueryVariable("year")
    if (year == false){
        year = "0"
    }
    var month_data = []
    $.ajax({
                type: "GET",//方法类型
                contentType: "application/json",
                url: "http://192.168.35.193:8081/admin/statistics/month/"+year,
                success: function (result) {
                    var lineChart = $('#line-chart');
                    if (lineChart.length > 0) {
                        new Chart(lineChart, {
                            type: 'line',
                            data: {
                                labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
                                datasets: [{
                                    label: 'Articles',
                                    data: result.arry,
                                    backgroundColor: 'rgba(66, 165, 245, 0.5)',
                                    borderColor: '#2196F3',
                                    borderWidth: 1
                                }]
                            },
                            options: {
                                legend: {
                                    display: false
                                },
                                scales: {
                                    yAxes: [{
                                        ticks: {
                                            beginAtZero: true
                                        }
                                    }]
                                }
                            }
                        });
                    }

                },
                beforeSend: function(xhr) {
                  xhr.setRequestHeader("X-Access-Token", token);
              },
                error : function() {
                  alert(result.error)
                }
            });

    //var lineChart = $('#line-chart');

    /*if (lineChart.length > 0) {
        new Chart(lineChart, {
            type: 'line',
            data: {
                labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
                datasets: [{
                    label: 'Users',
                    data: [parseInt(month_data[0]),0,0,0,0,0,0,0,0,0,0,0],
                    backgroundColor: 'rgba(66, 165, 245, 0.5)',
                    borderColor: '#2196F3',
                    borderWidth: 1
                }]
            },
            options: {
                legend: {
                    display: false
                },
                scales: {
                    yAxes: [{
                        ticks: {
                            beginAtZero: true
                        }
                    }]
                }
            }
        });
    }*/

    /**
     * Bar Chart
     */
     /*$.ajax({
                type: "GET",//方法类型
                contentType: "application/json",
                url: "http://192.168.35.193:8081/admin/statistics/month/"+year,
                success: function (result) {
                   

                },
                beforeSend: function(xhr) {
                  xhr.setRequestHeader("X-Access-Token", token);
              },
                error : function() {
                  alert(result.error)
                }
            });*/


    //柱状统计图
    $.ajax({
                type: "GET",//方法类型
                contentType: "application/json",
                url: "http://192.168.35.193:8081/admin/statistics/sort",
                success: function (result) {
                   var barChart = $('#bar-chart');
                   if (barChart.length > 0) {
                        new Chart(barChart, {
                            type: 'bar',
                            data: {
                                labels: result.sorts,
                                datasets: [{
                                    label: '# of Votes',
                                    data: result.arry,
                                    backgroundColor: [
                                        'rgba(244, 88, 70, 0.5)',
                                        'rgba(33, 150, 243, 0.5)',
                                        'rgba(0, 188, 212, 0.5)',
                                        'rgba(42, 185, 127, 0.5)',
                                        'rgba(156, 39, 176, 0.5)',
                                        'rgba(253, 178, 68, 0.5)'
                                    ],
                                    borderColor: [
                                        '#F45846',
                                        '#2196F3',
                                        '#00BCD4',
                                        '#2ab97f',
                                        '#9C27B0',
                                        '#fdb244'
                                    ],
                                    borderWidth: 1
                                }]
                            },
                            options: {
                                legend: {
                                    display: false
                                },
                                scales: {
                                    yAxes: [{
                                        ticks: {
                                            beginAtZero: true
                                        }
                                    }]
                                }
                            }
                        });
                    }

                },
                beforeSend: function(xhr) {
                  xhr.setRequestHeader("X-Access-Token", token);
              },
                error : function() {
                  alert(result.error)
                }
            });

    

    

    /**
     * Radar Chart
     */

     $.ajax({
                type: "GET",//方法类型
                contentType: "application/json",
                url: "http://192.168.35.193:8081/admin/statistics/gender",
                success: function (result) {
                    var radarChart = $('#radar-chart');

                    if (radarChart.length > 0) {
                        new Chart(radarChart, {
                            type: 'radar',
                            data: {
                                labels: ["08:00-12:00", "12:00-16:00", "16:00-20:00", "20:00-24:00", "24:00-04:00", "04:00-08:00"],
                                datasets: [{
                                    label: 'Male',
                                    data: result.male,
                                    backgroundColor: 'rgba(244, 88, 70, 0.5)',
                                    borderColor: '#F45846',
                                    borderWidth: 1
                                }, {
                                    label: 'Female',
                                    data: result.female,
                                    backgroundColor: 'rgba(33, 150, 243, 0.5)',
                                    borderColor: '#2196F3',
                                    borderWidth: 1
                                }]
                            }
                        });
                     }
                },
                beforeSend: function(xhr) {
                  xhr.setRequestHeader("X-Access-Token", token);
              },
                error : function() {
                  alert(result.error)
                }
            });



   
    

    /**
     * Pie Chart
     */
     $.ajax({
                type: "GET",//方法类型
                contentType: "application/json",
                url: "http://192.168.35.193:8081/admin/statistics/area",
                success: function (result) {
                   var pieChart = $('#pie-chart');
                   if (pieChart.length > 0) {
                        new Chart(pieChart, {
                            type: 'pie',
                            data: {
                                labels: result.area,
                                datasets: [{
                                    label: 'Users',
                                    data: result.array,
                                    backgroundColor: [
                                        'rgba(244, 88, 70, 0.5)',
                                        'rgba(33, 150, 243, 0.5)',
                                        'rgba(0, 188, 212, 0.5)',
                                        'rgba(42, 185, 127, 0.5)',
                                        'rgba(156, 39, 176, 0.5)',
                                        'rgba(253, 178, 68, 0.5)'
                                    ],
                                    borderColor: [
                                        'rgba(244, 88, 70, 0.5)',
                                        'rgba(33, 150, 243, 0.5)',
                                        'rgba(0, 188, 212, 0.5)',
                                        'rgba(42, 185, 127, 0.5)',
                                        'rgba(156, 39, 176, 0.5)',
                                        'rgba(253, 178, 68, 0.5)'
                                    ],
                                    borderWidth: 1
                                }]
                            }
                        });
    }

                },
                beforeSend: function(xhr) {
                  xhr.setRequestHeader("X-Access-Token", token);
              },
                error : function() {
                  alert(result.error)
                }
            });


    

    /**
     * Widget Line Chart
     */
    var wLineChart = $('.widget-line-chart');

    wLineChart.each(function (index, canvas) {
        new Chart(canvas, {
            type: 'line',
            data: {
                labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
                datasets: [{
                    label: 'Users',
                    data: [12, 19, 3, 5, 2, 3, 20, 33, 23, 12, 33, 10],
                    borderColor: '#fff',
                    borderWidth: 1,
                    fill: false,
                }]
            },
            options: {
                legend: {
                    display: false
                },
                scales: {
                    yAxes: [{
                        ticks: {
                            beginAtZero: true,
                            display: false,
                        },
                        gridLines: {
                            display: false,
                            drawBorder: false,
                        }
                    }],
                    xAxes: [{
                        ticks: {
                            display: false,
                        },
                        gridLines: {
                            display: false,
                            drawBorder: false,
                        }
                    }]
                }
            }
        });
    });
});
