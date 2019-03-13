# Test ssh connection

ansible -i ./hosts.cfg --user pi --ask-pass raspberry -m ping

# Deploy

ansible-playbook -i ./hosts.cfg --ask-pass playbook.yml
