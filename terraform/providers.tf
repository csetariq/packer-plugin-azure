terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>3.77"
    }
  }
}

provider "azurerm" {
  features {}
}

