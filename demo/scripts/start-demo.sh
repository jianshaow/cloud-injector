. sourceme.sh

echo start demo...

echo 1. show product resources
read -p 'press enter...'
./1.show-resources.sh
echo
read -p 'press enter...'
echo show pod details:
./show-pod.sh
h_yellow 'access the demo application in browser!'
read -p 'press enter...'
echo

echo 2. prepare a injection config
read -p 'press enter...'
./2.show-injection-config.sh
read -p 'press enter...'
echo

echo 3. label the demo namespace with injection=enabled and restart demo application
read -p 'press enter...'
./3.label-demo-ns.sh
read -p 'press enter...'
./restart-demo-app.sh
./wait-demo-app.sh
echo
read -p 'press enter...'
echo show pod details:
./show-pod.sh
h_yellow 'access the demo application in browser again! surprise!'
read -p 'press enter...'
echo

echo 4. remove the label injection=enabled from demo namespace and restart demo application
read -p 'press enter...'
./4.unlabel-demo-ns.sh
read -p 'press enter...'
./restart-demo-app.sh
./wait-demo-app.sh
echo
read -p 'press enter...'
echo show pod details:
./show-pod.sh
h_yellow 'access the demo application in browser again! everything is back to what it was!'
read -p 'press enter...'
echo

h_yellow 'Thank you!'
