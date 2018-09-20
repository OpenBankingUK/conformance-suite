defmodule Compliance.AccountsTest do
  @moduledoc """
  Tests for accounts context.
  """
  use Compliance.DataCase

  alias Compliance.Accounts

  import Mock
  import ExUnit.CaptureLog

  describe "users" do
    alias Compliance.Accounts.{User, UserValidationRun}

    @valid_attrs user_valid_attrs()
    @update_attrs %{
      email: @valid_attrs.email,
      first_name: "some updated first_name",
      last_name: "some updated last_name",
      provider: "some updated provider",
      token: "some updated token"
    }
    @invalid_attrs %{email: nil, first_name: nil, last_name: nil, provider: nil, token: nil}

    def assert_id_match(%User{} = user, %User{} = expected_user) do
      assert user.id == expected_user.id
      user
    end

    def assert_fields_match(%User{} = user, attrs) do
      assert user.email == attrs.email
      assert user.first_name == attrs.first_name
      assert user.last_name == attrs.last_name
      assert user.provider == attrs.provider
      assert user.token == attrs.token
      user
    end

    test "create_user_from_id_token" do
      google_client_id = "GOOGLE_CLIENT_ID"
      System.put_env("GOOGLE_OAUTH_CLIENT_ID", google_client_id)
      id_token =  "TOKEN"
      google_tokeninfo_url = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=#{id_token}"
      body = %{
        sub: @valid_attrs.token,
        aud: google_client_id,
        email: @valid_attrs.email,
        given_name: @valid_attrs.first_name,
        family_name: @valid_attrs.last_name
      }

      with_mock HTTPoison,
        get: fn ^google_tokeninfo_url ->
          {:ok, %HTTPoison.Response{status_code: 200, body: Poison.encode!(body)}}
        end do

        capture_log(fn ->
          assert {:ok, %User{} = user} = Accounts.create_user_from_id_token(id_token)
          user |> assert_fields_match(@valid_attrs)
        end)
      end

      with_mock HTTPoison,
        get: fn ^google_tokeninfo_url ->
          {:ok, %HTTPoison.Response{status_code: 400}}
        end do

        capture_log(fn ->
          assert {:error, "Error validating google_token_id"} = Accounts.create_user_from_id_token(id_token)
        end)
      end

      with_mock HTTPoison,
        get: fn ^google_tokeninfo_url ->
          {:ok, %HTTPoison.Response{status_code: 200, body: Poison.encode!(body |> Map.put("aud", "different"))}}
        end do
          capture_log(fn ->

          assert {:error, "Error: aud doesn't match GOOGLE_CLIENT_ID"}  = Accounts.create_user_from_id_token(id_token)
        end)
      end
    end

    test "list_users/0 returns all users" do
      existing_user = user_fixture()
      existing_user2 = user_fixture(email: "another_email", token: "other token")

      assert [%User{} = user, %User{} = user2] = Accounts.list_users()
      assert_id_match(user, existing_user)
      assert_id_match(user2, existing_user2)
    end

    test "get_user!/1 returns the user with given id" do
      existing_user = user_fixture()
      assert %User{} = user = Accounts.get_user!(existing_user.id)
      assert_id_match(user, existing_user)
    end

    test "get_user/1 returns the user with given token" do
      existing_user = user_fixture()
      token = existing_user.token
      assert %User{} = user = Accounts.get_user(token: token)
      assert_id_match(user, existing_user)
    end

    test "get_user/1 returns nil when no match for given token" do
      user_fixture()
      assert nil == Accounts.get_user(token: "bad_token")
    end

    test "create_user/1 with valid data creates a user" do
      assert {:ok, %User{} = user} = Accounts.create_user(@valid_attrs)
      user |> assert_fields_match(@valid_attrs)
    end

    test "create_user/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Accounts.create_user(@invalid_attrs)
    end

    test "create_or_update_user/1 with valid data creates a user" do
      assert {:ok, %User{} = user} = Accounts.create_or_update_user(@valid_attrs)
      user |> assert_fields_match(@valid_attrs)
    end

    test "create_or_update_user/1 with existing user email returns updated user" do
      assert {:ok, %User{} = existing_user} = Accounts.create_user(@valid_attrs)
      assert {:ok, %User{} = user} = Accounts.create_or_update_user(@update_attrs)

      user
      |> assert_id_match(existing_user)
      |> assert_fields_match(@update_attrs)
    end

    test "update_user/2 with valid data updates the user" do
      existing_user = user_fixture()
      assert {:ok, %User{} = user} = Accounts.update_user(existing_user, @update_attrs)

      user
      |> assert_id_match(existing_user)
      |> assert_fields_match(@update_attrs)
    end

    test "update_user/2 with invalid data returns error changeset and leaves persisted data unchanged" do
      existing_user = user_fixture()
      assert {:error, %Ecto.Changeset{}} = Accounts.update_user(existing_user, @invalid_attrs)
      user = Accounts.get_user!(existing_user.id)

      user
      |> assert_id_match(existing_user)
      |> assert_fields_match(@valid_attrs)
    end

    test "delete_user/1 deletes the user" do
      existing_user = user_fixture()
      assert {:ok, %User{} = user} = Accounts.delete_user(existing_user)
      user |> assert_id_match(existing_user)
      assert_raise Ecto.NoResultsError, fn -> Accounts.get_user!(existing_user.id) end
    end

    test "change_user/1 returns a user changeset" do
      user = user_fixture()
      assert %Ecto.Changeset{} = Accounts.change_user(user)
    end

    test "create_user_validation_run/2 creates user_validation_run" do
      user = user_fixture()
      validation_run_id = "test-val-run-id"
      user_id = user.id

      assert {:ok,
              %UserValidationRun{
                validation_run_id: ^validation_run_id,
                user_id: ^user_id,
                inserted_at: %NaiveDateTime{}
              }} = Accounts.create_user_validation_run(user, validation_run_id)
    end

    test "list_user_validation_runs/1 returns empty list for new user" do
      user = user_fixture()
      assert [] = Accounts.list_user_validation_runs(user)
    end

    test "list_user_validation_runs/1 returns list of user validation runs for given user" do
      user = user_fixture()
      validation_run_id = "test-val-run-id"
      validation_run_id2 = "test-val-run-id2"
      Accounts.create_user_validation_run(user, validation_run_id)
      Accounts.create_user_validation_run(user, validation_run_id2)

      other_user = user_fixture(email: "other email", token: "other token")
      validation_run_id3 = "test-val-run-id3"
      Accounts.create_user_validation_run(other_user, validation_run_id3)

      assert [
               %UserValidationRun{
                 validation_run_id: ^validation_run_id
               },
               %UserValidationRun{
                 validation_run_id: ^validation_run_id2
               }
             ] = Accounts.list_user_validation_runs(user)

      assert [
               %UserValidationRun{
                 validation_run_id: ^validation_run_id3
               }
             ] = Accounts.list_user_validation_runs(other_user)
    end

    test "has_user_validation_run?/2 returns true when run belongs to user, otherwise returns false" do
      user = user_fixture()
      validation_run_id = "test-val-run-id"
      Accounts.create_user_validation_run(user, validation_run_id)

      other_user = user_fixture(email: "other email", token: "other token")
      other_validation_run_id = "test-val-run-id2"
      Accounts.create_user_validation_run(other_user, other_validation_run_id)

      assert Accounts.has_user_validation_run?(user, validation_run_id)

      refute Accounts.has_user_validation_run?(other_user, validation_run_id)

      assert Accounts.has_user_validation_run?(other_user, other_validation_run_id)
    end
  end
end
