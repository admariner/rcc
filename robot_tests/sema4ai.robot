*** Settings ***
Library  OperatingSystem
Library  supporting.py
Resource  resources.robot
Suite Setup  Sema4.ai setup
Suite Teardown  Sema4.ai teardown
Default tags   WIP

*** Keywords ***
Sema4.ai setup
  Remove Directory  tmp/sema4home  True
  Prepare Sema4.ai Home    tmp/sema4home

Sema4.ai teardown
  Remove Directory  tmp/sema4home  True
  Prepare Robocorp Home    tmp/robocorp

*** Test cases ***

Goal: See rcc toplevel help for Sema4.ai
  Step        build/rcc --sema4ai --controller citests --help
  Must Have   SEMA4AI
  Must Have   Robocorp
  Must Have   --robocorp
  Must Have   --sema4ai
  Must Have   completion
  Must Have   robot
  Wont Have   ROBOCORP
  Wont Have   Robot
  Wont Have   assistant
  Wont Have   interactive
  Wont Have   community
  Wont Have   tutorial
  Wont Have   bash
  Wont Have   fish

Goal: See rcc commands for Sema4.ai
  Step        build/rcc --sema4ai --controller citests
  Use STDERR
  Must Have   SEMA4AI
  Wont Have   ROBOCORP
  Wont Have   Robocorp
  Wont Have   Robot
  Must Have   robot
  Wont Have   assistant
  Wont Have   interactive
  Wont Have   community
  Wont Have   tutorial
  Wont Have   completion
  Wont Have   bash
  Wont Have   fish

Goal: Default settings.yaml for Sema4.ai
  Step        build/rcc --sema4ai configuration settings --defaults --controller citests
  Must Have   Sema4.ai default settings.yaml
  Wont Have   assistant
  Wont Have   branding
  Wont Have   logo

Goal: Create package.yaml environment using uv
  Step        build/rcc --sema4ai ht vars -s sema4ai --controller citests robot_tests/bare_action/package.yaml
  Must Have   RCC_ENVIRONMENT_HASH=
  Must Have   RCC_INSTALLATION_ID=
  Must Have   SEMA4AI_HOME=
  Wont Have   ROBOCORP_HOME=
  Must Have   _4e67cd8_81359368
  Use STDERR
  Must Have   Progress: 01/15
  Must Have   Progress: 15/15
  Must Have   Running uv install phase.
