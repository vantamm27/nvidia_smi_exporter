#! /bin/bash

wget -O /usr/local/bin/nvidia_smi_exporter  https://github.com/vantamm27/nvidia_smi_exporter/raw/master/nvidia_smi_exporter

chmod +x /usr/local/bin/nvidia_smi_exporter

wget -O nvidia_smi_exporter.service https://raw.githubusercontent.com/vantamm27/nvidia_smi_exporter/master/nvidia_smi_exporter.service

useradd -m -s /bin/bash node_exporter

systemctl daemon-reload
systemctl enable nvidia_smi_exporter
systemctl start nvidia_smi_exporter
systemctl status nvidia_smi_exporter

mkdir -p /opt/iot/nvidia_smi_exporter


export HOSTNAME=`hostname`
export ENDPOINT="http://61.28.233.46:9091"

echo '#!/bin/bash' > test.sh 
echo "PUSHGATEWAY_SERVER=$ENDPOINT" >>  /opt/iot/nvidia_smi_exporter/nvidia_smi_exporter_metrics.sh
echo "NODE_NAME=$HOSTNAME" >> /opt/iot/nvidia_smi_exporter/nvidia_smi_exporter_metrics.sh 
echo "curl -s localhost:9101/metrics | curl --data-binary @- \$PUSHGATEWAY_SERVER/metrics/job/node-exporter/instance/\$NODE_NAME " >>  /opt/iot/nvidia_smi_exporter/nvidia_smi_exporter_metrics.sh 

chmod +x /opt/iot/nvidia_smi_exporter/nvidia_smi_exporter_metrics.sh

crontab -l | { cat; echo "*/1 * * * * /opt/iot/nvidia_smi_exporter/nvidia_smi_exporter_metrics.sh"; } | crontab -
