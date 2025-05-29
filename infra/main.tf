terraform {
    required_providers {
        mgc = {
            source = "magalucloud/mgc"
            version = "0.33.0"
        }
    }
}

provider "mgc" {
    api_key = ""
    region = "br-se1"
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
    key = ""
}

locals {
    primary_interface_id = [
        for interface in mgc_virtual_machine_instances.maquinas-test-devices-api.network_interfaces :
        interface.id if interface.primary
    ][0]
}

resource "mgc_network_public_ips_attach" "example" {
    public_ip_id = mgc_network_public_ips.meu_id_publico.id
    interface_id = local.primary_interface_id
}

resource "mgc_network_public_ips" "meu_id_publico" {
    description = "example public ip"
    vpc_id = "" // pegar pelo script apos a criação
}

resource "mgc_network_security_groups_attach" "regras_network"{
    security_group_id = ""
    interface_id = mgc_network_vpcs_interfaces.interface_gp.id

}

resource "mgc_network_vpcs_interfaces" "interface_gp" {
    name = "interface nomeada"
    vpc_id = ""
}