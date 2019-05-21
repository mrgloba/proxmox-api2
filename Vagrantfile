# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  config.vm.box = "cwiggs/proxmox5"

  config.vm.box_check_update = true

  config.vm.network "forwarded_port", guest: 8006, host: 8006

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "1024"
  end

  config.vm.provision "shell", inline: <<-SHELL
    sudo -E echo "deb http://download.proxmox.com/debian/pve stretch pve-no-subscription" > /etc/apt/sources.list.d/pve-enterprise.list
    export DEBIAN_FRONTEND=noninteractive
    sudo -E apt-get -qy update
    sudo -E apt-get -qy -o "Dpkg::Options::=--force-confdef" -o "Dpkg::Options::=--force-confold" dist-upgrade
    sudo -E apt-get -qy autoclean
	sudo -E pvesh create /access/users --userid testuser@pve --password testuser
   	sudo -E pvesh set /access/acl --path / --users testuser@pve --roles Administrator --propagate 1
  SHELL

  config.vm.provision :reload

  config.vm.provision "shell", inline: <<-SHELL
     sudo -E pvesh create /nodes/pve/aplinfo --storage local --template debian-9.0-standard_9.0-2_amd64.tar.gz
  SHELL

  
end
