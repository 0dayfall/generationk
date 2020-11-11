package main

import (
	"fmt"
	"os"
	"strconv"
)

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func createGraph(m *Context) {
	f, err := os.Create("graph.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(`<html>
	<head>
	  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
	  <script type="text/javascript">
		google.charts.load('current', {'packages':['corechart']});
		google.charts.setOnLoadCallback(drawChart);
  
	function drawChart() {
	  var data = google.visualization.arrayToDataTable([`)

	for _, ohlc := range m.asset.ohlc {
		f.WriteString("['" + ohlc.time.String() + "', " + floatToString(ohlc.open) + "," + floatToString(ohlc.high) + "," + floatToString(ohlc.low) + "," + floatToString(ohlc.close) + "],")
	}
	//fmt.Println(l, "bytes written successfully")

	l, _ = f.WriteString(`    ], true);

    var options = {
	  legend:'none',
	  'width':2400,
      'height':1800
    };

    var chart = new google.visualization.CandlestickChart(document.getElementById('chart_div'));

    chart.draw(data, options);
  }
    </script>
  </head>
  <body>
    <div id="chart_div" style="width: 900px; height: 500px;"></div>
  </body>
</html>`)

	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	fmt.Println(l, "bytes written successfully")
	err = f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
}
