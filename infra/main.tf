terraform {
    required_providers {
        mgc = {
            source = "magalucloud/mgc"
            version = "0.33.0"
        }
    }
}

provider "mgc" {
    api_key = var.mgc_api_key
    region = var.mgc_region
}

resource "mgc_virtual_machine_instances" "maquinas-test-devices-api" {
    name = "teste-devices-main"
    machine_type = "BV1-2-10"
    image = "cloud-ubuntu-24.04 LTS"
    ssh_key_name= mgc_ssh_keys.ssh_chave.name
    user_data = base64encode(file("arquivo-start.sh"))
}

resource "mgc_ssh_keys" "ssh_chave" {
    name = "infra-key"
    key = var.mgc_ssh_key
}

locals {
    primary_interface_id = [
        for interface in mgc_virtual_machine_instances.maquinas-test-devices-api.network_interfaces :
        interface.id if interface.primary
    ][0]
}

resource "mgc_network_public_ips_attach" "vm-public-ip-attachment" {
    public_ip_id = mgc_network_public_ips.meu_id_publico.id
    interface_id = local.primary_interface_id
}

resource "mgc_network_public_ips" "vm-public-ips" {
    description = "public ips"
    vpc_id = var.mgc_vpc_id 
}

resource "mgc_network_security_groups_attach" "regras_network"{
    security_group_id = var.mgc_sg_id
    interface_id = mgc_network_vpcs_interfaces.interface_gp.id

}

resource "mgc_network_vpcs_interfaces" "interface_gp" {
    name = "interface nomeada"
    vpc_id = var.mgc_vpc_i_id
}